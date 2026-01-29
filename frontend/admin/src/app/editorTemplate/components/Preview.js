'use client';
import { NewModal } from '@repo/shared-methods';
import { Button, Spin } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { EyeOutlined } from '@ant-design/icons';

// 预览样式
export const PreviewModal = ({ setOpen, value }) => {
	const [loading, setLoading] = useState(true);
	const [isShow, setIsShow] = useState(false);
	// 使用 useMemo 缓存 iframe 内容
	const iframeContent = useMemo(() => value, [value]);
	useEffect(() => {
		// 因为iframe导致弹框点击了打开很卡顿，所以加个定时器增加loading效果缓解
		let loadingTime = setTimeout(() => {
			setIsShow(true);
		}, 100);
		return () => {
			clearTimeout(loadingTime);
		};
	}, [iframeContent]);
	return (
		<NewModal width={850} title="Preview" onCancel={() => setOpen(false)} footer={null}>
			<div className="w-[100%] h-[500px] overflow-hidden">
				{loading && (
					<div className="flex justify-center items-center h-full">
						<Spin spinning={loading} />
					</div>
				)}
				{isShow && (
					<iframe
						className="w-[800px] h-[500px] overflow-hidden"
						srcDoc={iframeContent}
						sandbox="allow-scripts"
						// 添加性能优化属性
						loading="lazy"
						fetchPriority="low"
						onLoad={() => {
							setLoading(false);
						}}
					/>
				)}
			</div>
		</NewModal>
	);
};

export const PreviewBtn = ({ value, ...rest }) => {
	const [open, setOpen] = useState(false);
	return (
		<>
			<Button
				icon={<EyeOutlined />}
				type="link"
				size="small"
				onClick={() => setOpen(true)}
				{...rest}
			>
				Preview
			</Button>
			{open && <PreviewModal value={value} setOpen={setOpen} />}
		</>
	);
};
