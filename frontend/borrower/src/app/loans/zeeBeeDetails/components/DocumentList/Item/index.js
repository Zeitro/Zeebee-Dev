import React, { useState } from 'react';
import { Collapse, List, Space } from 'antd';
import {
	DeleteOutlined,
	FileImageOutlined,
	FilePdfOutlined,
	FileTextOutlined,
} from '@ant-design/icons';
import styles from '../index.module.scss';
import { DeleteBtn, ViewFile } from '@repo/shared-methods';
import { createRequest, eventBus } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const DocumentItem = ({ columns, item, isHide }) => {
	const [activeKey, setActiveKey] = useState([]);
	const [loading, setLoading] = useState(false);
	const { request } = createRequest();

	useDeferredSyncEffect({
		deferred: () => {
			setActiveKey(item.files?.length ? [item.id] : []);
		},
		deps: [item?.files?.length, item.id],
		delay: 0,
	});

	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/file/${id}`, {
				method: 'DELETE',
			});
			setLoading(false);
			eventBus.emit('getLoanDetail');
		} catch (error) {
			setLoading(false);
		}
	};
	const handleClickOpt = (key, record) => {
		if (key === 'delete') {
			handleDelete(record.id);
		}
	};

	const getIcon = (suffix) => {
		switch (suffix) {
			case 'pdf':
				return <FilePdfOutlined />;
			case 'png':
			case 'jpg':
			case 'jpeg':
				return <FileImageOutlined />;
			case 'doc':
			case 'docx':
				return <FileTextOutlined />;
			default:
				return <FileTextOutlined />;
		}
	};
	return (
		<Collapse
			size="default"
			collapsible="header"
			className={[styles.collapse, '!rounded-md overflow-hidden'].join(' ')}
			activeKey={activeKey}
			onChange={(key) => {
				setActiveKey(key);
			}}
			items={[
				{
					classNames: {
						header: 'bg-[#f5f5f5] !rounded-md !py-2 !px-4 w-full',
					},
					showArrow: false,
					key: item.id,
					label: (
						<div className="flex gap-4 items-center">
							{columns.map((columnsItem, index) => (
								<div
									key={index}
									style={{
										width: columnsItem.width,
										minWidth: columnsItem.minWidth,
										flex: columnsItem.minWidth ? 1 : '0 0 auto',
									}}
									className={[
										columnsItem?.ellipsis && styles.ellipsis,
										index === columns.length - 1 && 'flex items-center',
										columnsItem?.align === 'right' && 'justify-end',
										columnsItem?.align === 'center' && 'justify-center',
										!index && 'font-medium',
									].join(' ')}
								>
									{/* {item?.children?.length > 0 && !index && (
										<UpOutlined
											onClick={() => setExpanded(!expanded)}
											style={{
												transform: expanded
													? 'rotate(180deg)'
													: 'rotate(0deg)',
												transition: 'transform 0.3s',
												fontSize: '12px',
												marginRight: '4px',
											}}
										/>
									)} */}

									{columnsItem?.render?.(item[columnsItem.dataIndex], item) ||
										item[columnsItem.dataIndex] ||
										'--'}
								</div>
							))}
						</div>
					),
					children: (
						<List
							itemLayout="horizontal"
							dataSource={item?.files}
							renderItem={(fileItem, index) => {
								const name = fileItem.file_url.split('/').pop();
								const suffix = name.split('.').pop()?.toLowerCase();
								return (
									<List.Item key={index}>
										<div className="h-8 px-4 bg-white border border-base-border rounded-md flex items-center justify-center">
											{getIcon(suffix)}
											<div className="ml-1">{name}</div>
										</div>
										<Space>
											<ViewFile type="link" file={fileItem} />
											{!isHide && (
												<DeleteBtn
													icon={<DeleteOutlined />}
													onOk={() => handleClickOpt('delete', fileItem)}
													confirmLoading={loading}
												>
													Delete
												</DeleteBtn>
											)}
										</Space>
									</List.Item>
								);
							}}
						/>
					),
				},
			]}
		/>
	);
};

export default DocumentItem;
