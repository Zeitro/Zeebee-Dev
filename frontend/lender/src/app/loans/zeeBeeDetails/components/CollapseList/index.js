import React, { Children, useEffect, useState } from 'react';
import { Checkbox, Collapse, Form } from 'antd';
import styles from './index.module.scss';
const CollapseList = ({ list = [], initialValues = {} }) => {
	const [items, setItems] = useState([]);
	const [data, setData] = useState({});

	const formatList = (list) => {
		const newList = list?.map?.((item) => {
			if (item.children) {
				item.children = (
					<div>
						{item.children.map?.((child, index) => {
							return (
								<div
									key={index}
									className="flex justify-between items-center w-full"
								>
									<div>{child.label}</div>
									<Checkbox
										onChange={(e) =>
											setData({ ...data, [child.key]: e.target.checked })
										}
									/>
								</div>
							);
						})}
					</div>
				);
			} else {
				item.collapsible = 'icon';
			}
			item.extra = (
				<Checkbox onChange={(e) => setData({ ...data, [item.key]: e.target.checked })} />
			);
			return item;
		});
		return newList;
	};
	useEffect(() => {
		setItems(formatList(list));
	}, [list]);
	return (
		<Collapse
			accordion
			bordered
			ghost
			items={items}
			expandIconPosition="end"
			className={styles.collapseCommon}
		/>
	);
};
export default CollapseList;
