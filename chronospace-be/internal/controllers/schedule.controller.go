package controllers

import "chronospace-be/internal/services"

type ScheduleController struct {
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) *ScheduleController {
	return &ScheduleController{
		scheduleService: scheduleService,
	}
}