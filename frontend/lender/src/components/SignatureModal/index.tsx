'use client';
import type {ComponentRef} from 'react';
import {useRef} from 'react';
import SignaturePad from 'react-signature-pad-wrapper';
import {NewModal} from '@repo/shared-methods';
import {message} from 'antd';
//使用 ComponentRef 和 typeof 来获取组件的实例类型
type SignaturePadType = ComponentRef<typeof SignaturePad>;

interface SignatureModalProps {
	isOpen: boolean;
	onClose: () => void;
	// eslint-disable-next-line no-unused-vars
	onConfirm: (signature: string) => void;
}

export default function SignatureModal({ isOpen, onClose, onConfirm }: SignatureModalProps) {
	// ✅ 2. 使用新库提供的精确类型
	const sigCanvas = useRef<SignaturePadType>(null);
	const [messageApi, contextHolder] = message.useMessage();

	const handleClear = () => {
		sigCanvas.current?.clear();
	};

	const handleConfirm = () => {
		if (sigCanvas.current?.isEmpty()) {
			messageApi.open({
				type: 'warning',
				content: 'Please sign first!',
			});
			return;
		}

		// .toDataURL() 直接支持裁剪空白
		const signatureImage = sigCanvas.current?.toDataURL('image/png', 1.0);

		if (signatureImage) {
			onConfirm(signatureImage);
		}
	};

	if (!isOpen) {
		return null;
	}

	return (
		<>
			{contextHolder}
			<NewModal
				open={isOpen}
				title="Please sign in the space below："
				onCancel={onClose}
				cancelText="Clear"
				okText="Confirm"
				onOk={handleConfirm}
				cancelButtonProps={{
					onClick: handleClear,
				}}
				centered
				width={450}
			>
				<div className="w-full border border-base-border rounded-sm bg-gray-light">
					<SignaturePad
						ref={sigCanvas}
						options={{
							penColor: 'rgba(0, 0, 0, 0.88)',
							minWidth: 1.5, // 最小线宽（默认0.5，增大以加粗）
							maxWidth: 2.5, // 最大线宽（进一步加粗）
							throttle: 16, // 绘制频率（避免延迟导致的断线）
						}}
						canvasProps={{
							width: '400px',
							height: '200px',
							className: 'filter-contrast-120',
						}}
					/>
				</div>
			</NewModal>
		</>
	);
}
