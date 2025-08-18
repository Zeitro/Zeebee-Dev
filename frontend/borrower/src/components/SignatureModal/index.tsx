'use client';
import { useRef } from 'react';
import type { ComponentRef } from 'react';
import SignaturePad from 'react-signature-pad-wrapper';
import NewModal from '../NewModal';
//使用 ComponentRef 和 typeof 来获取组件的实例类型
type SignaturePadType = ComponentRef<typeof SignaturePad>;

interface SignatureModalProps {
	isOpen: boolean;
	onClose: () => void;
	onConfirm: (signature: string) => void;
}

export default function SignatureModal({ isOpen, onClose, onConfirm }: SignatureModalProps) {
	// ✅ 2. 使用新库提供的精确类型
	const sigCanvas = useRef<SignaturePadType>(null);

	if (!isOpen) {
		return null;
	}

	const handleClear = () => {
		sigCanvas.current?.clear();
	};

	const handleConfirm = () => {
		if (sigCanvas.current?.isEmpty()) {
			alert('请先签名');
			return;
		}

		// .toDataURL() 直接支持裁剪空白
		const signatureImage = sigCanvas.current?.toDataURL('image/png', 1.0);

		if (signatureImage) {
			onConfirm(signatureImage);
		}
	};

	return (
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
					options={{ penColor: 'black' }}
					canvasProps={{
						width: '400px',
						height: '200px',
						className: 'sigCanvas',
					}}
				/>
			</div>
		</NewModal>
	);
}
