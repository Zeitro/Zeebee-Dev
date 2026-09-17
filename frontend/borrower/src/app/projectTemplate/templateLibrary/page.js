'use client';

import { ProjectTemplateLibrary } from '@repo/project-management';
import { PERMISSIONS } from '@config/permissions';

export default function TemplateLibraryPage() {
	return <ProjectTemplateLibrary PERMISSIONS={PERMISSIONS} />;
}
