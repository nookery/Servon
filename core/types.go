package core

import (
	"servon/core/internal/contract"
	"servon/core/internal/libs"
	githubModels "servon/core/internal/libs/github/models"
	"servon/core/internal/models"
	"servon/core/internal/utils"
)

type OSType = libs.OSType
type CommandOptions = utils.CommandOptions
type CronTask = libs.CronTask
type ValidationError = libs.ValidationError
type ValidationErrors = libs.ValidationErrors
type Task = models.Task
type SoftwareInfo = contract.SoftwareInfo
type SuperSoft = contract.SuperSoft
type DeployLog = models.DeployLog
type WebhookPayload = githubModels.WebhookPayload
