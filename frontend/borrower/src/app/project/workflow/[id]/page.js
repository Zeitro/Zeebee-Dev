'use client';

import { WorkflowEditor } from '@repo/project-management';
import { useParams, useSearchParams } from 'next/navigation';

export default function ProjectWorkflowEditPage() {
	const params = useParams();
	const workflowId = params?.id;
	const searchParams = useSearchParams();
	const projectIdRaw = searchParams?.get('projectId');
	const projectId = projectIdRaw ? Number(projectIdRaw) : undefined;

	return <WorkflowEditor workflowId={workflowId} projectId={projectId} />;
}
