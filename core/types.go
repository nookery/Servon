package core

import (
	"servon/components/command_util"
	"servon/components/cron_util"
	"servon/components/github"
	"servon/components/logger"
	"servon/core/contract"
	"servon/core/managers"
	"servon/core/models"
)

type CommandOptions = command_util.CommandOptions
type CronTask = cron_util.CronTask
type ValidationError = cron_util.ValidationError
type ValidationErrors = cron_util.ValidationErrors
type Task = models.Task

type SoftwareInfo = contract.SoftwareInfo
type SuperSoft = contract.SuperSoft
type SuperGateway = contract.SuperGateway
type SuperService = contract.SuperService

type DeployLog = models.DeployLog
type WebhookPayload = github.WebhookPayload
type OSType = managers.OSType
type Project = contract.Project

type Deployer = contract.SuperDeployer

type LogUtil = logger.LogUtil
