'use client';

import { ProjectList } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';

export default function ProjectListPage() {
	return <ProjectList PERMISSIONS={PERMISSIONS} />;
}
