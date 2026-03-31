import React, { useEffect, useState } from 'react';
import { NewModal } from '@repo/shared-methods';
import Item from './Item';
const ListModal = ({ setOpen, list }) => {
	return (
		<NewModal width={572} title="Maturity Alerts" onCancel={() => setOpen(false)} footer={null}>
			<div className="flex flex-col gap-6">
				{list?.map((item, index) => (
					<Item key={index} item={item} />
				))}
			</div>
		</NewModal>
	);
};

export default ListModal;
