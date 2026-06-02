import React from 'react';
import { NewModal } from '@repo/shared-methods';
import { Button, Form, Input } from 'antd';
import useData from '../useData';

// 公共使用组件，用于在 loan 和 transaction 中创建文档类型
const AddRequest = ({
	open,
	setOpen,
	loan_id,
	category,
	id,
	apiType = 'loan', // loan 或 transaction
	title = 'Add request document type',
	onSuccessCallback = () => {},
}) => {
	const { form, loading, switchCreateDocumentType } = useData({
		id,
		category,
		setOpen,
		loan_id,
		apiType,
		onSuccessCallback,
	});

	return (
		<NewModal
			title={title}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={switchCreateDocumentType}
			okText="Save"
			confirmLoading={loading}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				initialValues={{
					file_types: [''],
				}}
				className="h-[400px] overflow-hidden overflow-y-auto"
			>
				<Form.List name="file_types">
					{(fields, { add, remove }) => {
						return (
							<>
								{fields?.map((field, index) => (
									<div
										key={field.key}
										className="flex justify-between items-start border-t [&:nth-of-type(1)]:border-t-0 border-base-border pt-4 first:pt-0"
									>
										<Form.Item
											label="Document type(Required)"
											name={[field.name]}
											rules={[
												{
													required: true,
													message: 'Please enter document type!',
												},
											]}
										>
											<Input
												placeholder="Enter document type"
												style={{ width: 300 }}
											/>
										</Form.Item>

										{fields.length > 1 && (
											<Button
												type="primary"
												ghost
												size="middle"
												onClick={() => remove(field.name)}
											>
												Delete
											</Button>
										)}
									</div>
								))}
								<Button size="middle" onClick={() => add()} type="primary" ghost>
									Add more document type
								</Button>
							</>
						);
					}}
				</Form.List>
			</Form>
		</NewModal>
	);
};

export default AddRequest;
