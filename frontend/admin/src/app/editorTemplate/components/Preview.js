'use client';
import { NewModal } from '@repo/shared-methods';
import { Button } from 'antd';
import { useState } from 'react';
import { EyeOutlined } from '@ant-design/icons';

// 预览样式
export const PreviewModal = ({ setOpen, value }) => {
	return (
		<NewModal width={850} title="Preview" onCancel={() => setOpen(false)} footer={null}>
			<div className="w-[100%] h-[500px] overflow-hidden">
				<iframe
					className="w-[800px] h-[500px] overflow-hidden"
					srcDoc={value}
					sandbox="allow-scripts"
				/>
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
