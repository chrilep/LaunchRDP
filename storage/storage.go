package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chrilep/LaunchRDP/config"
	"github.com/chrilep/LaunchRDP/models"
)

const (
	UsersFileName = "users.json"
	HostsFileName = "hosts.json"
)

// Storage handles reading and writing of users and hosts
type Storage struct {
	usersPath string
	hostsPath string
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{
		usersPath: config.GetConfigPath(UsersFileName),
		hostsPath: config.GetConfigPath(HostsFileName),
	}
}

// LoadUsers loads users from JSON file
func (s *Storage) LoadUsers() ([]models.User, error) {
	if _, err := os.Stat(s.usersPath); os.IsNotExist(err) {
		// File doesn't exist, return empty slice
		return []models.User{}, nil
	}

	data, err := os.ReadFile(s.usersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read users file: %w", err)
	}

	var users models.Users
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users.Users, nil
}

// SaveUsers saves users to JSON file (sorted alphabetically by username)
func (s *Storage) SaveUsers(users []models.User) error {
	// Sort users before saving to keep config file organized
	sortedUsers := make([]models.User, len(users))
	copy(sortedUsers, users)
	sort.Slice(sortedUsers, func(i, j int) bool {
		return strings.ToLower(sortedUsers[i].Username) < strings.ToLower(sortedUsers[j].Username)
	})

	usersData := models.Users{Users: sortedUsers}

	data, err := json.MarshalIndent(usersData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(s.usersPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

// LoadHosts loads hosts from JSON file
func (s *Storage) LoadHosts() ([]models.Host, error) {
	if _, err := os.Stat(s.hostsPath); os.IsNotExist(err) {
		// File doesn't exist, return empty slice
		return []models.Host{}, nil
	}

	data, err := os.ReadFile(s.hostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}

	var hosts models.Hosts
	if err := json.Unmarshal(data, &hosts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hosts: %w", err)
	}

	return hosts.Hosts, nil
}

// SaveHosts saves hosts to JSON file (sorted alphabetically by name)
func (s *Storage) SaveHosts(hosts []models.Host) error {
	// Sort hosts before saving to keep config file organized
	sortedHosts := make([]models.Host, len(hosts))
	copy(sortedHosts, hosts)
	sort.Slice(sortedHosts, func(i, j int) bool {
		return strings.ToLower(sortedHosts[i].Name) < strings.ToLower(sortedHosts[j].Name)
	})

	hostsData := models.Hosts{Hosts: sortedHosts}

	data, err := json.MarshalIndent(hostsData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal hosts: %w", err)
	}

	if err := os.WriteFile(s.hostsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write hosts file: %w", err)
	}

	return nil
}
