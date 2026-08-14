import React from 'react';
import Image from 'next/image';
import NodataSvg from '@repo/platform-core/assets/images/nodata.svg';
export const NoData = ({
	url,
	text = 'No Data',
	width,
	height,
	className,
	imageProps = {},
	children,
}) => {
	return (
		<div
			className={`h-full w-full flex flex-col justify-center items-center pb-5 pt-5 ${className}`}
			style={{ width, height }}
		>
			{children ? (
				children
			) : (
				<>
					{url ? (
						<Image src={url} alt="" width={120} height={120} {...imageProps} />
					) : (
						<NodataSvg style={{ width: 120, height: 120 }} {...imageProps} />
					)}
					<div className="text-base-label">{text}</div>
				</>
			)}
		</div>
	);
};
