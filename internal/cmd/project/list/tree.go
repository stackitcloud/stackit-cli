package list

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"

	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

type node struct {
	resourceID     string
	name           string
	lifecycleState resourcemanager.LifecycleState
	labels         map[string]string

	typ      resourceType
	parent   *node
	children []*node
}

type resourceType string

const (
	resourceTypeOrg     resourceType = "organization"
	resourceTypeFolder  resourceType = "folder"
	resourceTypeProject resourceType = "project"
)

type resourceTree struct {
	mu sync.Mutex

	authClient     *authorization.APIClient
	resourceClient *resourcemanager.APIClient
	member         string

	projectLifecycleState *string

	roots map[string]*node
}

func newResourceTree(resourceClient *resourcemanager.APIClient, authClient *authorization.APIClient, model *inputModel) (*resourceTree, error) {
	var member string
	if model.Member == nil {
		var err error
		member, err = auth.GetAuthEmail()
		if err != nil {
			return nil, fmt.Errorf("get email of authenticated user: %w", err)
		}
	} else {
		member = *model.Member
	}
	tree := &resourceTree{
		member:         member,
		resourceClient: resourceClient,
		authClient:     authClient,
		roots:          map[string]*node{},
	}
	if model.LifecycleState != "" {
		tree.projectLifecycleState = &model.LifecycleState
	}
	return tree, nil
}

func (r *resourceTree) Fill(ctx context.Context) error {
	resp, err := r.authClient.ListUserMemberships(ctx, r.member).ResourceType("organization").Execute()
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, orgMembership := range resp.GetItems() {
		g.Go(func() error {
			org, err := r.resourceClient.GetOrganizationExecute(ctx, orgMembership.GetResourceId())
			if err != nil {
				return err
			}
			orgNode := &node{
				resourceID:     org.GetOrganizationId(),
				name:           org.GetName(),
				typ:            resourceTypeOrg,
				lifecycleState: org.GetLifecycleState(),
				labels:         org.GetLabels(),
			}
			r.mu.Lock()
			r.roots[orgNode.resourceID] = orgNode
			r.mu.Unlock()
			if err := r.fillNode(ctx, orgNode); err != nil {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

func (r *resourceTree) fillNode(ctx context.Context, parent *node) error {
	if err := r.getNodeProjects(ctx, parent); err != nil {
		return err
	}
	req := r.resourceClient.ListFolders(ctx).ContainerParentId(parent.resourceID)
	resp, err := req.Execute()
	if err != nil {
		if !isForbiddenError(err) {
			return err
		}
		// listing folder for parent was forbidden, trying with member
		resp, err = req.Member(r.member).Execute()
		if err != nil {
			return err
		}
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, folder := range resp.GetItems() {
		g.Go(func() error {
			newFolderNode := &node{
				resourceID:     folder.GetFolderId(),
				parent:         parent,
				typ:            resourceTypeFolder,
				name:           folder.GetName(),
				lifecycleState: resourcemanager.LIFECYCLESTATE_ACTIVE,
				labels:         folder.GetLabels(),
			}
			parent.children = append(parent.children, newFolderNode)
			return r.fillNode(ctx, newFolderNode)
		})
	}
	return g.Wait()
}

func (r *resourceTree) getNodeProjects(ctx context.Context, parent *node) error {
	req := r.resourceClient.ListProjects(ctx).ContainerParentId(parent.resourceID)
	resp, err := req.Execute()
	if err != nil {
		if !isForbiddenError(err) {
			return err
		}
		// listing projects for parent was forbidden, trying with member
		resp, err = req.Member(r.member).Execute()
		if err != nil {
			return err
		}
	}
	for _, proj := range resp.GetItems() {
		if r.projectLifecycleState != nil && *r.projectLifecycleState != strings.ToLower(string(proj.GetLifecycleState())) {
			continue
		}
		projNode := &node{
			resourceID:     proj.GetProjectId(),
			typ:            resourceTypeProject,
			name:           proj.GetName(),
			labels:         proj.GetLabels(),
			lifecycleState: proj.GetLifecycleState(),
			parent:         parent,
		}
		parent.children = append(parent.children, projNode)
	}
	return nil
}
