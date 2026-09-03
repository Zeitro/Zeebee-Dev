'use client';

import { ProjectMyTask } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';
export default function ProjectMyTaskPage() {
	return <ProjectMyTask PERMISSIONS={PERMISSIONS} />;
}
