package core

import (
	"servon/core/internal/contract"
	"servon/core/internal/managers"
	"servon/core/internal/managers/deployers"
	"servon/core/internal/managers/github"
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
type SuperGateway = contract.SuperGateway
type SuperService = contract.SuperService

type DeployLog = models.DeployLog
type WebhookPayload = github.WebhookPayload
type OSType = managers.OSType
type Project = contract.Project

type Deployer = deployers.Deployer

type LogUtil = utils.LogUtil
