package phase

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/files"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

var (
	createDirectory  = files.CreateDirectory
	copyContents     = files.CopyContents
	ownershipMatches = files.OwnershipMatches
	chownRecursive   = files.ChownRecursive
)

// Prepare prepares the runtime
func Prepare(cfg config.Config, logger *zap.Logger) error {
	uid, gid, err := parseOwnership(cfg)
	if err != nil {
		return err
	}
	if cfg.AgentHome != "" {
		logger.Debug("creating agent home",
			zap.String("path", cfg.AgentHome),
		)
		if err := createDirectory(cfg.AgentHome); err != nil {
			return fmt.Errorf("create agent home %q: %w", cfg.AgentHome, err)
		}
		logger.Debug("agent home created",
			zap.String("path", cfg.AgentHome),
		)
		source := filepath.Join(cfg.RuntimeHome, "config", "agent")
		logger.Debug("copying shared agent configuration",
			zap.String("source", source),
			zap.String("destination", cfg.AgentHome),
		)
		if err := copyContents(source, cfg.AgentHome); err != nil {
			return fmt.Errorf("copy shared agent configuration from %q to %q: %w", source, cfg.AgentHome, err)
		}
		logger.Debug("shared agent configuration copied",
			zap.String("source", source),
			zap.String("destination", cfg.AgentHome),
		)
	} else {
		logger.Debug("agent home not configured; shared agent configuration skipped")
	}
	matches, err := ownershipMatches(cfg.UserHome, uid, gid)
	if err != nil {
		return fmt.Errorf("check ownership for %q: %w", cfg.UserHome, err)
	}
	if !matches {
		logger.Debug("updating home ownership",
			zap.String("path", cfg.UserHome),
			zap.Int("uid", uid),
			zap.Int("gid", gid),
		)
		if err := chownRecursive(cfg.UserHome, uid, gid); err != nil {
			return fmt.Errorf("update ownership for %q: %w", cfg.UserHome, err)
		}
		logger.Debug("home ownership updated",
			zap.String("path", cfg.UserHome),
			zap.Int("uid", uid),
			zap.Int("gid", gid),
		)
	} else {
		logger.Debug("home ownership already matches",
			zap.String("path", cfg.UserHome),
			zap.Int("uid", uid),
			zap.Int("gid", gid),
		)
	}
	return nil
}

func parseOwnership(cfg config.Config) (int, int, error) {
	uid, err := strconv.Atoi(cfg.LocalUID)
	if err != nil {
		return 0, 0, fmt.Errorf("local_uid must be a numeric user id: %w", err)
	}
	gid, err := strconv.Atoi(cfg.LocalGID)
	if err != nil {
		return 0, 0, fmt.Errorf("local_gid must be a numeric group id: %w", err)
	}
	return uid, gid, nil
}
