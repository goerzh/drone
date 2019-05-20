// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/drone/drone/cmd/drone-server/config"
	"github.com/drone/drone/handler/api"
	"github.com/drone/drone/handler/web"
	"github.com/drone/drone/livelog"
	"github.com/drone/drone/operator/manager"
	"github.com/drone/drone/pubsub"
	"github.com/drone/drone/service/commit"
	"github.com/drone/drone/service/hook/parser"
	"github.com/drone/drone/service/license"
	"github.com/drone/drone/service/org"
	"github.com/drone/drone/service/repo"
	"github.com/drone/drone/service/token"
	"github.com/drone/drone/service/user"
	"github.com/drone/drone/store/batch"
	config2 "github.com/drone/drone/store/config"
	"github.com/drone/drone/store/cron"
	"github.com/drone/drone/store/perm"
	"github.com/drone/drone/store/secret"
	"github.com/drone/drone/store/secret/global"
	"github.com/drone/drone/store/step"
	"github.com/drone/drone/trigger"
	cron2 "github.com/drone/drone/trigger/cron"
)

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Injectors from wire.go:

func InitializeApplication(config3 config.Config) (application, error) {
	client := provideClient(config3)
	refresher := provideRefresher(config3)
	db, err := provideDatabase(config3)
	if err != nil {
		return application{}, err
	}
	userStore := provideUserStore(db)
	renewer := token.Renewer(refresher, userStore)
	commitService := commit.New(client, renewer)
	cronStore := cron.New(db)
	repositoryStore := provideRepoStore(db)
	fileService := provideContentService(client, renewer)
	configService := provideConfigPlugin(client, fileService, config3)
	statusService := provideStatusService(client, renewer, config3)
	buildStore := provideBuildStore(db)
	stageStore := provideStageStore(db)
	scheduler := provideScheduler(stageStore, config3)
	system := provideSystem(config3)
	webhookSender := provideWebhookPlugin(config3, system)
	triggerer := trigger.New(configService, commitService, statusService, buildStore, scheduler, repositoryStore, userStore, webhookSender)
	cronScheduler := cron2.New(commitService, cronStore, repositoryStore, userStore, triggerer)
	coreLicense := provideLicense(client, config3)
	datadog := provideDatadog(userStore, repositoryStore, buildStore, system, coreLicense, config3)
	corePubsub := pubsub.New()
	logStore := provideLogStore(db, config3)
	logStream := livelog.New()
	netrcService := provideNetrcService(client, renewer, config3)
	encrypter, err := provideEncrypter(config3)
	if err != nil {
		return application{}, err
	}
	secretStore := secret.New(db, encrypter)
	globalSecretStore := global.New(db, encrypter)
	stepStore := step.New(db)
	buildManager := manager.New(buildStore, configService, corePubsub, logStore, logStream, netrcService, repositoryStore, scheduler, secretStore, globalSecretStore, statusService, stageStore, stepStore, system, userStore, webhookSender)
	secretService := provideSecretPlugin(config3)
	registryService := provideRegistryPlugin(config3)
	runner := provideRunner(buildManager, secretService, registryService, config3)
	hookService := provideHookService(client, renewer, config3)
	licenseService := license.NewService(userStore, repositoryStore, buildStore, coreLicense)
	permStore := perm.New(db)
	repositoryService := repo.New(client, renewer)
	session := provideSession(userStore, config3)
	batcher := batch.New(db)
	syncer := provideSyncer(repositoryService, repositoryStore, userStore, batcher, config3)
	server := api.New(buildStore, commitService, cronStore, corePubsub, globalSecretStore, hookService, logStore, coreLicense, licenseService, permStore, repositoryStore, repositoryService, scheduler, secretStore, stageStore, stepStore, statusService, session, logStream, syncer, system, triggerer, userStore, webhookSender)
	organizationService := orgs.New(client, renewer)
	userService := user.New(client)
	admissionService := provideAdmissionPlugin(client, organizationService, userService, config3)
	hookParser := parser.New(client)
	middleware := provideLogin(config3)
	options := provideServerOptions(config3)
	configStore := config2.New(db)
	webServer := web.New(admissionService, buildStore, client, hookParser, coreLicense, licenseService, middleware, repositoryStore, session, syncer, triggerer, userStore, userService, webhookSender, options, system, configStore)
	handler := provideRPC(buildManager, config3)
	metricServer := provideMetric(session, config3)
	mux := provideRouter(server, webServer, handler, metricServer)
	serverServer := provideServer(mux, config3)
	mainApplication := newApplication(cronScheduler, datadog, runner, serverServer, userStore)
	return mainApplication, nil
}
