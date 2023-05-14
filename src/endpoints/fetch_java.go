package endpoints

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"torch/src/structs"
	"torch/src/utils"

	"github.com/gin-gonic/gin"
)

func FetchJava(host string, port uint16) (*structs.JavaStatus, error) {
	originalHost, originalPort := host, port

	if port == 25565 {
		srv, _ := LookupSrv(host)
		if srv != nil {
			host, port = srv.Target, srv.Port
		}
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), statusTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(statusTimeout)); err != nil {
		return nil, err
	}

	if err = sendHandshake(conn, host, port); err != nil {
		return nil, err
	}

	if err = sendStatusRequest(conn); err != nil {
		return nil, err
	}

	rawJavaResponse, err := readStatusResponse(conn)
	if err != nil {
		return nil, err
	}

	pingStart := time.Now()
	pingPayload := time.Now().UnixNano()

	if err = sendPing(conn, pingPayload); err != nil {
		return nil, err
	}

	if err = readPong(conn, pingPayload); err != nil {
		return nil, err
	}

	return createJavaStatus(originalHost, originalPort, host, port, rawJavaResponse, pingStart), nil
}

func sendHandshake(conn net.Conn, host string, port uint16) error {
	buf := &bytes.Buffer{}

	if _, err := utils.WriteVarInt(0x00, buf); err != nil {
		return err
	}

	if _, err := utils.WriteVarInt(47, buf); err != nil {
		return err
	}

	if err := utils.WriteString(host, buf); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, port); err != nil {
		return err
	}

	if _, err := utils.WriteVarInt(1, buf); err != nil {
		return err
	}

	return utils.WritePacket(buf, conn)
}

func sendStatusRequest(conn net.Conn) error {
	buf := &bytes.Buffer{}

	if _, err := utils.WriteVarInt(0x00, buf); err != nil {
		return err
	}

	return utils.WritePacket(buf, conn)
}

func readStatusResponse(conn net.Conn) (structs.RawJavaStatus, error) {
	var rawJavaResponse structs.RawJavaStatus

	_, _, err := utils.ReadVarInt(conn)
	if err != nil {
		return rawJavaResponse, err
	}

	packetId, _, err := utils.ReadVarInt(conn)
	if err != nil {
		return rawJavaResponse, err
	}

	if packetId != 0x00 {
		return rawJavaResponse, fmt.Errorf("unexpected packet ID (expected 0x00, got 0x%02x)", packetId)
	}

	response, err := utils.ReadString(conn)
	if err != nil {
		return rawJavaResponse, err
	}

	err = json.Unmarshal([]byte(response), &rawJavaResponse)
	return rawJavaResponse, err
}

func sendPing(conn net.Conn, pingPayload int64) error {
	buf := &bytes.Buffer{}

	if _, err := utils.WriteVarInt(0x01, buf); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, pingPayload); err != nil {
		return err
	}

	return utils.WritePacket(buf, conn)
}

func readPong(conn net.Conn, pingPayload int64) error {
	_, _, err := utils.ReadVarInt(conn)
	if err != nil {
		return err
	}

	packetId, _, err := utils.ReadVarInt(conn)
	if err != nil {
		return err
	}

	if packetId != 0x01 {
		return fmt.Errorf("unexpected packet ID (expected 0x01, got 0x%02x)", packetId)
	}

	var returnedPayload int64
	if err := binary.Read(conn, binary.BigEndian, &returnedPayload); err != nil {
		return err
	}

	if returnedPayload != pingPayload {
		return fmt.Errorf("unexpected payload (expected %d, got %d)", pingPayload, returnedPayload)
	}

	return nil
}

func createJavaStatus(originalHost string, originalPort uint16, host string, port uint16, rawJavaResponse structs.RawJavaStatus, pingStart time.Time) *structs.JavaStatus {
	// Process data
	description := structs.Parse(rawJavaResponse.Description)

	samplePlayers := make([]structs.Player, 0)

	if rawJavaResponse.Players.Sample != nil {
		for _, player := range rawJavaResponse.Players.Sample {
			name := structs.Parse(player.Name)
			samplePlayers = append(samplePlayers, structs.Player{
				ID:   player.ID,
				Name: *name,
			})
		}
	}

	versionText := structs.Parse(rawJavaResponse.Version.Name)

	var srv *structs.SrvRecord
	if originalHost == host {
		srv = nil
	} else {
		srv = &structs.SrvRecord{
			Host: host,
			Port: port,
		}
	}

	if rawJavaResponse.Favicon == "" {
		rawJavaResponse.Favicon = defaultFavicon
	}

	result := &structs.JavaStatus{
		Host: originalHost,
		Port: originalPort,
		Version: structs.Version{
			Name:     versionText,
			Protocol: rawJavaResponse.Version.Protocol,
		},
		Players: structs.Players{
			Max:    rawJavaResponse.Players.Max,
			Online: rawJavaResponse.Players.Online,
			Sample: samplePlayers,
		},
		Description: description,
		Favicon:     rawJavaResponse.Favicon,
		SrvRecord:   srv,
		Latency:     time.Duration(time.Since(pingStart).Milliseconds()),
		ModInfo:     nil,
		ObtainedAt:  time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(statusCacheTime)),
	}

	if len(rawJavaResponse.ModInfo.Type) > 0 {
		mods := make([]structs.Mod, 0)
		for _, mod := range rawJavaResponse.ModInfo.List {
			mods = append(mods, structs.Mod{
				ID:      mod.ModID,
				Version: mod.Version,
			})
		}
		result.ModInfo = &structs.ModInfo{
			Type:    rawJavaResponse.ModInfo.Type,
			ModList: mods,
		}
	}

	if rawJavaResponse.ForgeData.Mods != nil {
		mods := make([]structs.Mod, 0)
		for _, mod := range rawJavaResponse.ForgeData.Mods {
			mods = append(mods, structs.Mod{
				ID:      mod.ModID,
				Version: mod.Version,
			})
		}
		result.ModInfo = &structs.ModInfo{
			Type:    "forge",
			ModList: mods,
		}
	}
	return result
}

func FetchJavaHandler(c *gin.Context) {
	host := c.Query("host")
	port := c.Query("port")

	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 25565
	}

	uintPort := uint16(intPort)

	cacheKey := fmt.Sprintf("%s:%d", host, intPort)
	data, err := javaCache.Value(cacheKey)
	if err == nil {
		c.JSON(200, data.Data().(*structs.JavaStatus))
		return
	}

	fetchedData, err := FetchJava(host, uintPort)
	if err != nil {
		c.JSON(200, structs.OfflineServer{
			Offline: true,
			Host:    host,
			Port:    uintPort,
		})
		return
	}

	javaCache.Add(cacheKey, statusCacheTime, fetchedData)
	c.JSON(200, fetchedData)
}
