'use client';

import { ProjectTemplateMine } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';

export default function MyTemplatePage() {
	return <ProjectTemplateMine PERMISSIONS={PERMISSIONS} />;
}
