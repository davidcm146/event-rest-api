package main

import (
	"net/http"

	"github.com/davidcm146/event-rest-api/internal/database"
	"github.com/gin-gonic/gin"
)

// createEvent creates a new event
//
// @Summary Create a new event
// @Description Create a new event
// @Tags events
// @Accept json
// @Produce json
// @Param event body database.Event true "Event to create"
// @Success 201 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events [post]
// @Security BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Events.Insert(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// getAllEvents return all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving events"})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found"})
		return
	}
	c.JSON(http.StatusOK, &events)
}

// getEvent returns an event by ID
//
// @Summary Get event by ID
// @Description Get a single event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} database.Event
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
	}
	c.JSON(http.StatusOK, event)
}

// @Summary Update an event
// @Description Update an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body database.Event true "Updated event data"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [put]
// @Security BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id := c.Param("id")

	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}
	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent deletes an event
//
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 204 {object} nil
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id} [delete]
// @Security BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id := c.Param("id")
	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// addAttendeeToEvent adds a user as attendee to event
//
// @Summary Add attendee to event
// @Description Add a user as attendee to an event
// @Tags attendees
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param userId path string true "User ID to add as attendee"
// @Success 201 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id}/attendees/{userId} [post]
// @Security BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId := c.Param("id")
	userId := c.Param("userId")
	user := app.getUserFromContext(c)

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add attendees to this event"})
		return
	}

	userToAdd, err := app.models.Users.GetById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing attendee"})
		return
	}

	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already an attendee of this event"})
		return
	}

	attendee := &database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding attendee to event"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User added to event successfully", "attendee": attendee})
}

// getAttendeesByEvent returns list of attendees for event
//
// @Summary Get attendees for event
// @Description Get list of users attending an event
// @Tags attendees
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {array} database.User
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesByEvent(c *gin.Context) {
	eventId := c.Param("id")

	users, err := app.models.Attendees.GetAttendeesByEventId(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving attendees for event"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// removeAttendeeFromEvent removes an attendee from event
//
// @Summary Remove attendee from event
// @Description Remove a user from attendees of an event
// @Tags attendees
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param userId path string true "User ID to remove"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/events/{id}/attendees/{userId} [delete]
// @Security BearerAuth
func (app *application) removeAttendeeFromEvent(c *gin.Context) {
	eventId := c.Param("id")
	userId := c.Param("userId")
	user := app.getUserFromContext(c)
	event, err := app.models.Events.Get(eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	attendee, err := app.models.Attendees.GetByEventAndAttendee(eventId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving attendee"})
		return
	}

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add attendees to this event"})
		return
	}

	if attendee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
		return
	}

	if err := app.models.Attendees.Delete(userId, eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing attendee from event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendee removed from event successfully"})
}

// getEventsByAttendee returns events an attendee is participating
//
// @Summary Get events by attendee
// @Description Get events that a user is attending
// @Tags attendees
// @Accept json
// @Produce json
// @Param id path string true "Attendee ID"
// @Success 200 {array} database.Event
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id := c.Param("id")

	events, err := app.models.Attendees.GetEventsByAttendeeId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving events for attendee"})
		return
	}
	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for this attendee"})
		return
	}
	c.JSON(http.StatusOK, events)
}
