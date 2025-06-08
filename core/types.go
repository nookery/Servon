package core

import (
	"servon/components/command_util"
	"servon/components/github"
	"servon/components/log_util"
	"servon/core/internal/contract"
	"servon/core/internal/managers"
	"servon/core/internal/managers/deployers"
	"servon/core/internal/models"
)

type CommandOptions = command_util.CommandOptions
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

type LogUtil = log_util.LogUtil
