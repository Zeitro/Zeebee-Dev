'use client';

import { WorkflowEditor } from '@repo/project-management';
import { useParams } from 'next/navigation';

export default function ProjectWorkflowEditPage() {
	const params = useParams();
	const workflowId = params?.id;

	return <WorkflowEditor workflowId={workflowId} />;
}
