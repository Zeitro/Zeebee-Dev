'use client';

import { ProjectDetails } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';

export default function ProjectDetailPage() {
	return <ProjectDetails PERMISSIONS={PERMISSIONS} />;
}
