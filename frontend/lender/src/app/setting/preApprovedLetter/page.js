'use client';
import { Button, Form, message } from 'antd';
import { useEffect, useState } from 'react';
import { FieldTips } from '@repo/shared-methods';
import { createRequest, injectHeadResources, getBodyInnerHTML } from '@repo/shared-methods/utils';
import useContent from '../system/components/PreApprovalLetter/useContent';
import Configuration from './components/Configuration';
import useConfiguration from './components/Configuration/useConfiguration';
import DownloadPdf from './components/DownloadPdf';
import _ from 'lodash';
const PreApprovedLetter = () => {
	const { switchBtn } = useContent();
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [valueHtml, setValueHtml] = useState(''); // 模板内容（只有body内的html）
	const [completeTemplate, setCompleteTemplate] = useState(''); // 最新的模板数据(包含head/link/style/body)
	const [data, setData] = useState({}); //左侧配置项
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [open, setOpen] = useState(false);
	const [params, setParams] = useState({});
	const { updateVariablesInHtml } = useConfiguration({
		valueHtml,
		setValueHtml,
		data,
	});

	const handlePreview = async (result) => {
		const { content } = await request('/organization_template/preview', {
			method: 'POST',
			body: {
				...result,
			},
		});
		setCompleteTemplate(content);
		injectHeadResources(content);
		let bodyHtml = getBodyInnerHTML(content);
		let newHtml = bodyHtml;

		const exampleVariablesObj = JSON.parse(result?.example_variables);

		result?.variable_group_schema?.map((item) => {
			item.variables.map((i) => {
				const vKey = i.key.replace('.', '');
				const exampleValue = exampleVariablesObj[vKey] || `{${vKey}}`;
				let v = i.source === 'Dynamic' ? exampleValue : i.fixed_value;
				newHtml = updateVariablesInHtml(newHtml, v, i.key, i.format, false);
			});
		});
		setData(result);
		form.setFieldsValue(result);
		setValueHtml(newHtml);
	};

	const getDetails = async () => {
		try {
			setLoading(true);
			const params = new URLSearchParams({ purpose: 'preApprovedLetter' });
			const res = await request(`/organization_template/detail_by_purpose?${params}`);
			await handlePreview(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	const handleReset = async () => {
		try {
			setLoading(true);
			const res = await request(`/organization_template/${data?.id}/reset`, {
				method: 'POST',
			});
			await handlePreview(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			await request(`/organization_template/${data?.id}`, {
				method: 'PUT',
				body: {
					...values,
				},
			});
			message.destroy();
			message.open({
				type: 'success',
				content: 'Update successful！',
			});
			setConfirmLoading(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	useEffect(() => {
		const fetchDetails = async () => {
			await getDetails();
		};
		fetchDetails();
	}, []);

	return (
		<div className="flex-1 flex overflow-hidden bg-[#f5f5f5] mt-[2px]">
			<div className="border-r border-base-border w-[450px] flex-shrink-0 h-full bg-background flex flex-col relative px-8 pt-6">
				<h2 className="border-b-2 border-base-border text-xl font-medium pb-4 mb-4 flex items-center justify-between flex-shrink-0 text-black-dark">
					<span>
						Configuration
						<FieldTips
							title={
								<>
									For example data, please set your own pre-approved letter, which
									will be used for the pre-approved loan submitted by the
									borrower. Zeebee will automatically send this attachment to the
									borrower.
									<br />
									<br />
									Note: The "Default" option is calculated by the system, while
									"Customize" is set by yourself.
								</>
							}
							isShowDescription={true}
							rootClassName="!max-w-[300px]"
							className="ml-1"
						/>
					</span>
					<FieldTips title={'Reset data and get the latest configuration template.'}>
						<Button size="small" onClick={handleReset}>
							Reset
						</Button>
					</FieldTips>
				</h2>
				<Configuration
					valueHtml={valueHtml}
					setValueHtml={setValueHtml}
					data={data}
					configurationForm={form}
					loading={loading}
					confirmLoading={confirmLoading}
				/>
				<div className="w-full flex-shrink-0 bg-background pt-1 pb-4 flex flex-col gap-3">
					<Button
						type="primary"
						className="w-full"
						loading={confirmLoading}
						onClick={handleSave}
						size="large"
					>
						Save
					</Button>
					<Button
						className="w-full"
						onClick={() => {
							const v = form.getFieldsValue();
							v.variable_group_schema = v?.variable_group_schema?.map((item) => {
								item.variables = item.variables.map((i) => {
									i.source = 'Template';
									return i;
								});
								return item;
							});
							setParams(_.cloneDeep(v));
							setOpen(true);
						}}
						size="large"
					>
						Download
					</Button>
					<div className="flex items-center justify-between">
						<div className="text-xs text-base-label font-semibold">
							Send Pre-Approval Letter
							<FieldTips
								title="Note: When activated, your borrower will automatically receive the pre-approval letter configured on this page."
								isShowDescription={true}
								rootClassName="!max-w-[300px]"
								className="ml-1 text-xs"
							/>
						</div>
						{switchBtn}
					</div>
				</div>
			</div>

			<div
				className="flex-1 overflow-hidden overflow-y-auto py-6 flex justify-center items-center [&>.page]:max-w-[800px]"
				dangerouslySetInnerHTML={{ __html: valueHtml }}
			></div>
			{open && (
				<DownloadPdf
					setOpen={setOpen}
					valueHtml={valueHtml}
					data={data}
					inputSettingParams={params}
					completeTemplate={completeTemplate}
					updateVariablesInHtml={updateVariablesInHtml}
				/>
			)}
		</div>
	);
};
export default PreApprovedLetter;
