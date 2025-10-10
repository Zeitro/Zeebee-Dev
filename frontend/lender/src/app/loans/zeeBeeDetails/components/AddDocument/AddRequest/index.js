import React, { use } from 'react';
import NewModal from '@/components/NewModal';
import { Button, Form, Input, Select } from 'antd';
const AddRequest = ({ open, setOpen }) => {
	const [form] = Form.useForm();

	const handleOK = async () => {
		const values = await form.validateFields();
	};
	return (
		<NewModal
			title={'Add request document type'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Save"
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				initialValues={{
					document_type: [
						{
							type: '1',
						},
					],
				}}
				className="h-[400px] overflow-hidden overflow-y-auto"
			>
				<Form.List name="document_type">
					{(fields, { add, remove }) => {
						return (
							<>
								{fields?.map((field, index) => (
									<div
										key={field.key}
										className="flex justify-between items-start not-last:border-b border-base-border mb-4"
									>
										<Form.Item
											label="Document type(Required)"
											name={[field.name, 'type']}
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
