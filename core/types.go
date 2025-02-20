package core

import (
	"servon/core/internal/contract"
	githubModels "servon/core/internal/integrations/github/models"
	"servon/core/internal/libs/managers"
	"servon/core/internal/models"
	"servon/core/internal/utils"
)

type CommandOptions = utils.CommandOptions
type CronTask = managers.CronTask
type ValidationError = managers.ValidationError
type ValidationErrors = managers.ValidationErrors
type Task = models.Task
type SoftwareInfo = contract.SoftwareInfo
type SuperSoft = contract.SuperSoft
type DeployLog = models.DeployLog
type WebhookPayload = githubModels.WebhookPayload
type OSType = managers.OSType
