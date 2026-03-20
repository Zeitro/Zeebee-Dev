// 插入变量
export class CustomVariablesMenu {
	static variables = [];
	constructor() {
		this.title = 'Variable';
		this.tag = 'select';
		this.width = 80;
	}

	insertVariableFn = (editor, variableName) => {
		if (editor && variableName) {
			editor.focus();
			// 插入带样式的变量占位符
			editor.dangerouslyInsertHtml(
				`<span data-variable="${variableName}" style="background: #f6ffed; color: #389e0d; border-radius: 3px; contenteditable="false">{{ ${variableName} }}</span>`,
			);
		}
	};

	// 下拉框的选项
	getOptions(editor) {
		const options = [
			{ value: '', text: 'Variable', disabled: true },
			...CustomVariablesMenu.variables.map((variable) => ({
				value: variable.key,
				text: variable.label,
			})),
		];
		return options;
	}

	// 菜单是否需要激活
	isActive(editor) {
		return true;
	}

	// 获取菜单执行时的 value
	getValue(editor) {
		return ''; // 默认值
	}

	// 菜单是否需要禁用
	isDisabled(editor) {
		return false;
	}

	// 点击菜单时触发的函数
	exec(editor, value) {
		if (value && this.insertVariableFn) {
			this.insertVariableFn(editor, value);
		}
	}
	// 更新变量数组的方法
	updateVariables(newVariables) {
		CustomVariablesMenu.variables = newVariables;
	}
}

// 按钮菜单,弹出自定义modal框
export class CustomViewCodeMenu {
	constructor(showModal) {
		this.title = 'View Code';
		this.tag = 'button';
		this.showModal = showModal; // 保存方法
	}
	getValue(editor) {
		return '';
	}

	// 菜单是否需要激活（如选中加粗文本，“加粗”菜单会激活），用不到则返回 false
	isActive(editor) {
		return false;
	}

	// 菜单是否需要禁用（如选中 H1 ，“引用”菜单被禁用），用不到则返回 false
	isDisabled(editor) {
		return false;
	}

	exec(editor, value) {
		if (this.isDisabled(editor)) return;
		if (this.showModal) this.showModal(); // 触发弹窗显示
	}
}
