package routes

import (
	"backend/broker"
	"backend/database"
	"backend/prisma/db"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetPublicTasks(c *fiber.Ctx) error {
	ctx := context.Background()
	tasks, err := database.Client().Task.FindMany(
		db.Task.Public.Equals(true),
	).With(db.Task.User.Fetch()).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot get public tasks"})
	}
	return c.JSON(tasks)
}

func GetTasks(c *fiber.Ctx) error {
	userID := int(c.Locals("userID").(float64))

	ctx := context.Background()
	tasks, err := database.Client().Task.FindMany(
		db.Task.UserID.Equals(userID),
	).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot get tasks"})
	}
	return c.JSON(tasks)
}

func CreateTask(c *fiber.Ctx, brk *broker.Broker) error {
	userID := int(c.Locals("userID").(float64))

	var input struct {
		Title   string `json:"title" validate:"required,min=3,max=30"`
		Content string `json:"content" validate:"required,min=3,max=250"`
		Public  bool   `json:"public" validate:"boolean"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	json_validate := validator.New()
	if err := json_validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	task, err := database.Client().Task.CreateOne(
		db.Task.Title.Set(input.Title),
		db.Task.User.Link(db.User.ID.Equals(userID)),
		db.Task.Content.Set(input.Content),
		db.Task.Public.Set(input.Public),
	).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create task"})
	}

	// Broadcast create event to all clients
	if input.Public {
		brk.Broadcast(broker.TaskEvent{
			Type:   broker.EventCreate,
			TaskID: task.ID,
			Data:   task,
		})
	}

	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx, brk *broker.Broker) error {
	userID := int(c.Locals("userID").(float64))
	taskID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	var input struct {
		Title     *string `json:"title" validate:"min=3,max=30"`
		Content   *string `json:"content" validate:"required,min=3,max=250"`
		Completed *bool   `json:"completed" validate:"boolean"`
		Public    *bool   `json:"public" validate:"boolean"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	json_validate := validator.New()
	if err := json_validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := context.Background()
	task, err := database.Client().Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Exec(ctx)

	if err != nil || task.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	updateParams := []db.TaskSetParam{}
	if input.Title != nil {
		updateParams = append(updateParams, db.Task.Title.Set(*input.Title))
	}
	if input.Content != nil {
		updateParams = append(updateParams, db.Task.Content.Set(*input.Content))
	}
	if input.Completed != nil {
		updateParams = append(updateParams, db.Task.Completed.Set(*input.Completed))
	}
	if input.Public != nil {
		updateParams = append(updateParams, db.Task.Public.Set(*input.Public))
	}

	updatedTask, err := database.Client().Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Update(updateParams...).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update task"})
	}

	// Broadcast update event to all clients
	brk.Broadcast(broker.TaskEvent{
		Type:   broker.EventUpdate,
		TaskID: taskID,
		Data:   updatedTask,
	})

	return c.JSON(updatedTask)
}

func DeleteTask(c *fiber.Ctx, brk *broker.Broker) error {
	userID := int(c.Locals("userID").(float64))
	taskID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	ctx := context.Background()
	task, err := database.Client().Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Exec(ctx)

	if err != nil || task.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	_, err = database.Client().Task.FindUnique(
		db.Task.ID.Equals(taskID),
	).Delete().Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete task"})
	}

	// Broadcast delete event to all clients
	brk.Broadcast(broker.TaskEvent{
		Type:   broker.EventDelete,
		TaskID: taskID,
	})

	return c.SendStatus(fiber.StatusNoContent)
}
