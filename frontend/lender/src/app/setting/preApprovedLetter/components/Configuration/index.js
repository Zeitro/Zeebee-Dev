import { Form, Spin } from 'antd';
import { AssemblyFormItem } from '@repo/shared-methods';
import useConfiguration from './useConfiguration';
const Configuration = ({
	confirmLoading,
	configurationForm,
	valueHtml,
	setValueHtml,
	data,
	loading = false,
	name = 'configuration-form',
	className,
	isAllEditable,
}) => {
	const { handleFieldsChange, updateVariablesInHtml } = useConfiguration({
		valueHtml,
		setValueHtml,
		data,
	});

	return (
		<Form
			form={configurationForm}
			name={name}
			layout="vertical"
			autoComplete="off"
			className={[
				'flex-1 overflow-hidden overflow-y-auto [&_.ant-form-item]:mb-4 [&_.ant-form-item-explain-error]:leading-4',
				className,
			].join(' ')}
			size="large"
			disabled={confirmLoading}
			scrollToFirstError
			onFieldsChange={handleFieldsChange}
		>
			<Spin spinning={loading}>
				<Form.List name="variable_group_schema">
					{(fields) =>
						fields.map(({ key, name, ...restField }, index) => {
							const item = data?.variable_group_schema?.[index];
							return (
								<div key={key}>
									<div className="border-b border-base-border font-medium py-1 mb-4 text-base/[22px] text-black-dark">
										{item?.name}
									</div>
									<Form.List
										name={[name, 'variables']}
										{...restField}
										rules={[{ required: true }]}
									>
										{(fields2) =>
											fields2.map(
												(
													{
														key: variableItemKey,
														name: variableItemName,
														...variableItemRestField
													},
													i,
												) => {
													const variableItem = item?.variables?.[i];
													return variableItem?.is_editable ? (
														<AssemblyFormItem
															key={variableItemKey}
															item={variableItem}
															formListRestField={
																variableItemRestField
															}
															formListName={[variableItemName]}
															allFormListName={[
																'variable_group_schema',
																name,
																'variables',
																variableItemName,
															]}
															example_variables={
																data?.example_variables
															}
															changeHtml={(key, value, format) => {
																updateVariablesInHtml(
																	valueHtml,
																	value,
																	key,
																	format,
																);
															}}
															isAllEditable={isAllEditable}
														/>
													) : (
														''
													);
												},
											)
										}
									</Form.List>
								</div>
							);
						})
					}
				</Form.List>
			</Spin>
		</Form>
	);
};

export default Configuration;
