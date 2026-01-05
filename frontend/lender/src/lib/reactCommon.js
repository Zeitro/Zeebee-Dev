import React from 'react';

// 将svg内的fill修改为currentColor，方便使用
export function overrideFill(node, targetColor = '#原始颜色', newColor = 'currentColor') {
	if (!node) return node;
	console.log('99999', node.props);

	if (node.props && node.props.fill === targetColor) {
		console.log(node);

		return React.cloneElement(node, { fill: newColor });
	}
	if (node.props && node.props.children) {
		return React.cloneElement(node, {
			children: React.Children.map(node.props.children, (child) =>
				overrideFill(child, targetColor, newColor),
			),
		});
	}
	return node;
}
