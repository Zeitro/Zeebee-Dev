'use client';
import { Button, Image } from 'antd';
import React, { useState } from 'react';
import ViewSvg from '@/assets/images/view.svg';
import styles from './index.module.scss';
// 预览文件
const ViewFile = ({ children, file, ...rest }) => {
	const [visible, setVisible] = useState(false);
	const handleView = () => {
		const suffix = (file?.name?.split('.')?.[1] || '')?.toLowerCase();
		switch (suffix) {
			case 'jpg':
			case 'jpeg':
			case 'png':
				setVisible(true);
				break;
			default:
				window.open(
					'https://zeebee-files-staging.s3.us-west-1.amazonaws.com/files-dev/pdf/f96858d6-2c9c-400b-8b1b-079e53e1f6a2',
					'_blank',
				);
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
					src: file?.url,
					onVisibleChange: (vis) => setVisible(vis),
				}}
				alt=""
			/>
		</>
	);
};

export default ViewFile;
