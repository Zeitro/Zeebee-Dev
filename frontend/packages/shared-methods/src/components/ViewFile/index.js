'use client';
import { Button, Image } from 'antd';
import React, { useState } from 'react';
import ViewSvg from '@/assets/images/view.svg';
import styles from './index.module.scss';
// 预览文件
export const ViewFile = ({ children, file, fileKey = 'file_url', ...rest }) => {
	const [visible, setVisible] = useState(false);
	const handleView = () => {
		const suffix = (file?.[fileKey]?.split('.')?.pop() || '')?.toLowerCase();
		switch (suffix) {
			case 'jpg':
			case 'jpeg':
			case 'png':
				setVisible(true);
				break;
			default:
				window.open(file?.[fileKey], '_blank');
				break;
		}
	};
	return (
		<>
			<Button
				icon={<ViewSvg className={styles.viewIcon} />}
				size="large"
				onClick={handleView}
				{...rest}
			>
				{children || 'View'}
			</Button>
			<Image
				width={0}
				height={0}
				preview={{
					visible,
					src: file?.[fileKey],
					onVisibleChange: (vis) => setVisible(vis),
				}}
				alt=""
			/>
		</>
	);
};
