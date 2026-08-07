package global

import (
	"gorm.io/gorm"

	"zeebee-backend/common/email"
	"zeebee-backend/common/rbac"
	"zeebee-backend/common/redisx"
)

// This file collects test-only setters for the package-level singletons that
// borrower logic reads (Config / redis / db / enforcer / email). Production
// wires them via the unexported init* funcs in InitGlobal; tests build an
// isolated sqlite + miniredis environment and inject it here. Mirrors the
// lender / admin harness.

// SetRedis replaces the package redis client. Test-only.
func SetRedis(cli *redisx.Client) { redisDB = cli }

// SetGormDB replaces the package gorm DB. Test-only.
func SetGormDB(d *gorm.DB) { db = d }

// SetBorrowerEnforcer replaces the package Casbin enforcer. Test-only.
func SetBorrowerEnforcer(e *rbac.Enforcer) { BorrowerEnforcer = e }

// SetEmailClient replaces the package email client. Test-only.
func SetEmailClient(e *email.Email) { EmailClient = e }
