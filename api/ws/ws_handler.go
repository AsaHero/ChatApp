package ws

import (
	"github.com/AsaHero/chat_app/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub         *Hub
	userService service.User
}

func NewHandler(hub *Hub, userService service.User) *fiber.App {
	handler := &Handler{
		hub:         hub,
		userService: userService,
	}

	app := fiber.New()

	app.Post("/room", handler.CreateRoom)
	app.Get("/room", handler.GetRooms)
	app.Get("/client/:roomId", handler.GetClients)

	app.Use("/join/:id", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})
	app.All("/join/:id", websocket.New(handler.JoinRoom))

	return app
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"err": "cannot parse request body"})
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	return c.Status(fiber.StatusOK).JSON(req)
}

type GetRoomsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *fiber.Ctx) error {
	rooms := make([]GetRoomsResponse, 0)

	for _, v := range h.hub.Rooms {
		rooms = append(rooms, GetRoomsResponse{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return c.Status(fiber.StatusOK).JSON(rooms)
}

type GetClientsResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *fiber.Ctx) error {
	var clients []GetClientsResponse
	roomID := c.Params("roomId")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		clients = make([]GetClientsResponse, 0)
		return c.Status(fiber.StatusOK).JSON(clients)
	}

	for _, v := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, GetClientsResponse{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	return c.Status(fiber.StatusOK).JSON(clients)
}

func (h *Handler) JoinRoom(c *websocket.Conn) {
	roomId := c.Params("id")
	userId := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     c,
		ID:       userId,
		RoomID:   roomId,
		Username: username,
		Message:  make(chan *Message),
	}

	m := &Message{
		Content:  "new user has join the room",
		RoomID:   roomId,
		Username: username,
	}

	// register user from register channel
	h.hub.Register <- cl
	// broadcast message

	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}
