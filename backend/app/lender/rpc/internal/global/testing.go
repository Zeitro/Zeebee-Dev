package global

import (
	"gorm.io/gorm"

	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
	"zeebee-backend/domains/zeebee/limiter"
)

// This file collects test-only setters for the package-level singletons that
// lender logic reads (Config / redis / db / enforcer / email / limiter).
// Production wires them via the unexported init* funcs in InitGlobal; tests
// build an isolated sqlite + miniredis environment and inject it here.
// Mirrors the tenant testutil wiring pattern (svcCtx injection is not enough
// because lender logic reads these globals directly).

// SetRedis replaces the package redis client. Test-only.
func SetRedis(cli *redisx.Client) { redisDB = cli }

// SetGormDB replaces the package gorm DB. Test-only.
func SetGormDB(d *gorm.DB) { db = d }

// SetLenderEnforcer replaces the package Casbin enforcer. Test-only.
func SetLenderEnforcer(e *rbac.Enforcer) { LenderEnforcer = e }

// SetEmailClient replaces the package email client. Test-only.
func SetEmailClient(e *email.Email) { EmailClient = e }

// SetLimiterMgr replaces the package rate limiter manager. Test-only.
func SetLimiterMgr(m *limiter.LimiterManager) { LimiterMgr = m }
