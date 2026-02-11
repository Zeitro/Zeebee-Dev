import { NewModal, DownloadFileButton } from '@repo/shared-methods';
import { replaceBodyContent } from '@repo/shared-methods/utils';
import { useEffect, useMemo, useState } from 'react';
import Configuration from './Configuration';
import { Form } from 'antd';
import _ from 'lodash';

const DownloadPdf = ({
	setOpen,
	valueHtml,
	inputSettingParams,
	completeTemplate,
	updateVariablesInHtml,
}) => {
	const [downloadHtml, setDownloadHtml] = useState('');
	const [setting, setSetting] = useState({});
	const [form] = Form.useForm();
	const htmlV = useMemo(() => {
		return replaceBodyContent(completeTemplate, downloadHtml);
	}, [downloadHtml, completeTemplate]);

	useEffect(() => {
		const newData = _.cloneDeep(inputSettingParams);

		if (valueHtml) {
			let newHtml = valueHtml;
			newData?.variable_group_schema?.map((item) => {
				item.variables.map((i) => {
					newHtml = updateVariablesInHtml(newHtml, i.fixed_value, i.key, i.format, false);
				});
			});
			// 提取状态更新逻辑为独立函数
			const updateDownloadHtmlAndForm = (newHtml, newData) => {
				setDownloadHtml(newHtml);
				form.setFieldsValue(newData);
				setSetting(newData);
			};
			// 调用提取的函数进行状态更新
			updateDownloadHtmlAndForm(newHtml, newData);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [inputSettingParams, valueHtml]);
	return (
		<NewModal
			title="Download"
			width={'80vw'}
			onCancel={() => setOpen(false)}
			okText="Send request"
			footer={null}
			className="[&_.ant-modal-content_.ant-modal-body]:!px-0 [&_.ant-modal-content_.ant-modal-body]:!py-0"
		>
			<div className="flex h-[70vh]">
				<div className="flex flex-col w-1/3 max-w-[400px] bg-[#f5f5f5] px-6 py-4">
					<div className="text-base-label pb-3">
						Customize the letter details below for this specific download. These changes
						will not affect your saved template.
					</div>
					<Configuration
						valueHtml={downloadHtml}
						setValueHtml={setDownloadHtml}
						data={setting}
						configurationForm={form}
						name="downloadForm"
						isAllEditable={true}
						className="[&::-webkit-scrollbar-thumb]:border-[#f5f5f5] [&::-webkit-scrollbar]:bg-[#f5f5f5]"
					/>
				</div>
				<div
					className="flex-1 overflow-hidden overflow-y-auto py-6 flex justify-center items-center [&>.page]:max-w-[800px] bg-base-gray"
					dangerouslySetInnerHTML={{ __html: downloadHtml }}
				></div>
			</div>
			<div className="py-4 px-4">
				<DownloadFileButton
					type="primary"
					className="w-full"
					fileUrl="/pdf/html_to_pdf"
					requestOptions={{
						method: 'POST',
						isExternalLink: false,
						body: {
							content: htmlV,
							name: 'preApprovalLetter.pdf',
						},
					}}
				>
					Download PDF
				</DownloadFileButton>
			</div>
		</NewModal>
	);
};
export default DownloadPdf;
