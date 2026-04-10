import React, {useState} from 'react';
import {NewModal} from '@repo/shared-methods';
import {Button, Form, Input, message} from 'antd';
import {createRequest, eventBus} from '@repo/shared-methods/utils';

const AddRequest = ({ open, setOpen, loan_id, category }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);

	const handleOK = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			await request('/loan_document', {
				method: 'POST',
				body: {
					...values,
					loan_id,
					category,
				},
			});
			message.open({
				type: 'success',
				content: 'Add Document request successfully',
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<NewModal
			title={'Add request document type'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
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
