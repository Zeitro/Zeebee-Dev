'use client';
import React from 'react';
import dynamic from 'next/dynamic';

const ApplicationFormViewPDF = dynamic(
	() => import('@/components/ApplicationFormViewPDF'), // Adjust the path to your component
	{
		ssr: false, // This is the crucial part
		// Optional: Add a loading component to improve UX
		loading: () => <p>Loading...</p>,
	},
);

const PDF = () => {
	return (
		<div className="h-full w-full">
			<ApplicationFormViewPDF url="/api/getApplication" isFullLink={true} />
		</div>
	);
};
export default PDF;
