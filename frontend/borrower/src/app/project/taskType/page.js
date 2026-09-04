'use client';

import { ProjectTaskType } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';
export default function ProjectTaskTypePage() {
	return <ProjectTaskType PERMISSIONS={PERMISSIONS} />;
}
