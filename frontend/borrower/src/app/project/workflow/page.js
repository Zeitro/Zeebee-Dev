'use client';

import { ProjectWorkflow } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';
export default function ProjectWorkflowPage() {
	return <ProjectWorkflow PERMISSIONS={PERMISSIONS} />;
}
