import React from 'react';
const colorInfo = {
	Requested: {
		border: '1px solid #91D5FF',
		background: '#E6F7FF',
		color: 'var(--base-primary)',
	},
	Pending: {
		border: 'var(--base-border)',
		background: 'var(--background)',
		color: 'var(--foreground)',
	},
	'Ready for review': {
		border: '1px solid #B7EB8F',
		background: '#F6FFED',
		color: '#52C41A',
	},
};
const DocumentStatus = ({ status = 'Pending' }) => {
	return (
		<div style={{ ...colorInfo?.[status], borderRadius: 2 }} className="px-2 w-min">
			{status}
		</div>
	);
};

export default DocumentStatus;
