import React, { useCallback, useEffect, useState } from 'react';
import { Checkbox, Collapse } from 'antd';
import styles from './index.module.scss';

const CollapseList = ({ list = [], setList, hideCollapsible, isHideCheck = false }) => {
	const [items, setItems] = useState([]);

	const changeList = useCallback(
		(checked, index, childIndex) => {
			let newList = [...list];
			if (typeof childIndex === 'number') {
				newList[index].children[childIndex].checked = checked;
				const allCheck = newList[index].children.every((item) => item.checked);
				newList[index].checked = allCheck;
			} else {
				newList[index].checked = checked;
			}
			setList([...newList]);
		},
		[list, setList],
	);

	const formatList = useCallback(
		(list) => {
			const newList = list?.map?.((item, i) => {
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
										{!isHideCheck && (
											<Checkbox
												checked={child.checked}
												onChange={(e) =>
													changeList(e.target.checked, i, index)
												}
											/>
										)}
									</div>
								);
							})}
						</div>
					);
				} else {
					item.collapsible = 'icon';
				}
				item.extra = !isHideCheck ? (
					<Checkbox
						checked={item.checked}
						onChange={(e) => changeList(e.target.checked, i)}
					/>
				) : (
					''
				);
				return item;
			});
			return newList;
		},
		[changeList, isHideCheck],
	);
	useEffect(() => {
		setItems(formatList(list));
	}, [formatList, list]);
	return (
		<Collapse
			accordion
			bordered
			ghost
			items={items}
			expandIconPosition="end"
			className={[styles.collapseCommon, hideCollapsible && styles.collapseHide].join(' ')}
		/>
	);
};
export default CollapseList;
