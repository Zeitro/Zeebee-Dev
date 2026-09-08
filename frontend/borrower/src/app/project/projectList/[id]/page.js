'use client';

import { ProjectKanban } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';

export default function ProjectKanbanPage() {
	return <ProjectKanban PERMISSIONS={PERMISSIONS} />;
}
