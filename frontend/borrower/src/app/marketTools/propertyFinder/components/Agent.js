import { Popover, Typography } from 'antd';
import { useState } from 'react';
import { CheckCircleOutlined, CloseOutlined, CopyOutlined } from '@ant-design/icons';
import { formatPhone } from '@repo/shared-methods/utils';
const { Paragraph } = Typography;

export const getCopyContent = (content, options = {}) => {
	return (
		<Paragraph
			copyable={{
				icon: [<CopyOutlined key="copy-icon" />, <CheckCircleOutlined key="copied-icon" />],
				tooltips: ['click here', 'Copied!'],
				text: content,
			}}
			ellipsis={{
				rows: 1,
				tooltip: {
					title: content,
					placement: 'bottom',
				},
			}}
			{...options}
			className={['!mb-0', options.className].join(' ')}
		>
			{content}
		</Paragraph>
	);
};
const Agent = ({ record }) => {
	const [open, setOpen] = useState(false);

	return (
		<Popover
			open={open}
			onOpenChange={(open) => setOpen(open)}
			content={
				<div className="w-[280px]">
					<div className="flex justify-between items-center pb-2 border-b border-base-border">
						<span className="text-base font-medium">SELLER AGENT</span>
						<CloseOutlined className="cursor-pointer" onClick={() => setOpen(false)} />
					</div>
					{record.agent_name &&
						getCopyContent(record.agent_name, {
							className: 'text-base pt-1',
							strong: true,
						})}
					{record.agent_phone &&
						getCopyContent(formatPhone(record.agent_phone), {
							className: 'text-base pt-1',
						})}
					{record.agent_email &&
						getCopyContent(record.agent_email, {
							className: 'text-base pt-1',
						})}
				</div>
			}
			trigger="click"
			placement="right"
		>
			<span className="text-base-primary cursor-pointer">{record.agent_name}</span>
		</Popover>
	);
};
export default Agent;
