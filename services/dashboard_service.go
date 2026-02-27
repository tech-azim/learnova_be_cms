package services

import (
	"github.com/tech-azim/be-learnova/models"
	"github.com/tech-azim/be-learnova/repositories"
)

type DashboardService interface {
	GetDashboardData() (DashboardResponse, error)
}

type DashboardResponse struct {
	TotalProgram        int64                         `json:"total_program"`
	TotalRegistration   int64                         `json:"total_registration"`
	ActiveParticipants  int64                         `json:"active_participants"`
	PendingParticipants int64                         `json:"pending_participants"`
	LatestRegistrations []models.Registration         `json:"latest_registrations"`
	RecentActivities    []ActivityItem                `json:"recent_activities"`
	PopularPrograms     []repositories.PopularProgram `json:"popular_programs"`
}

type ActivityItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Program   string `json:"program"`
	Status    string `json:"status"`
	TimeAgo   string `json:"time_ago"`
	UpdatedAt string `json:"updated_at"`
}

type dashboardService struct {
	repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
	return &dashboardService{repo}
}

func (s *dashboardService) GetDashboardData() (DashboardResponse, error) {
	var response DashboardResponse

	totalProgram, err := s.repo.GetTotalProgram()
	if err != nil {
		return response, err
	}
	response.TotalProgram = totalProgram

	totalRegistration, err := s.repo.GetTotalRegistration()
	if err != nil {
		return response, err
	}
	response.TotalRegistration = totalRegistration

	activeParticipants, err := s.repo.GetActiveParticipants()
	if err != nil {
		return response, err
	}
	response.ActiveParticipants = activeParticipants

	pendingParticipants, err := s.repo.GetPendingParticipants()
	if err != nil {
		return response, err
	}
	response.PendingParticipants = pendingParticipants

	latestRegistrations, err := s.repo.GetLatestRegistrations(5)
	if err != nil {
		return response, err
	}
	response.LatestRegistrations = latestRegistrations

	recentActivities, err := s.repo.GetRecentActivities(10)
	if err != nil {
		return response, err
	}
	response.RecentActivities = mapToActivityItems(recentActivities)

	popularPrograms, err := s.repo.GetPopularPrograms(5)
	if err != nil {
		return response, err
	}
	response.PopularPrograms = popularPrograms

	return response, nil
}

func mapToActivityItems(registrations []models.Registration) []ActivityItem {
	items := make([]ActivityItem, 0, len(registrations))
	for _, r := range registrations {
		items = append(items, ActivityItem{
			ID:        r.ID,
			Name:      r.Name,
			Program:   r.Program.Title,
			Status:    r.Status,
			UpdatedAt: r.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return items
}
