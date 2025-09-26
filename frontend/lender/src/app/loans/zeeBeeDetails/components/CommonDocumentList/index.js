import React from 'react';
const CommonDocumentList = ({
	list = [
		{
			total: 0,
		},
	],
}) => {
	return (
		<div>
			{list.map((item, index) => (
				<div key={index}>{item.label}</div>
			))}
		</div>
	);
};
export default CommonDocumentList;
