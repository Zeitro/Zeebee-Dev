'use client';
import { useEffect, useRef, useState, forwardRef, useImperativeHandle } from 'react';
import { Input } from 'antd';
import DistanceSvg from '@repo/shared-methods/assets/images/distance.svg';
const GeoLocationInput = ({ placeholder, onLocationSelect, onError, ...props }, ref) => {
	const [inputValue, setInputValue] = useState('');
	const [isLoading, setIsLoading] = useState(false);
	const [isSelecting, setIsSelecting] = useState(false);
	const [errorTips, setErrorTips] = useState('');
	const autocompleteRef = useRef(null);
	const inputRef = useRef(null);

	// 选择了google下拉后，获取经纬度地址
	const handlePlaceSelect = () => {
		setIsSelecting(true);
		setIsLoading(true);
		setErrorTips('');
		try {
			const place = autocompleteRef.current.getPlace();

			if (!place.geometry) {
				setErrorTips('No location data available for this input.');
				setIsLoading(false);
				return;
			}

			const location = {
				lat: place.geometry.location.lat(),
				lng: place.geometry.location.lng(),
				address: place.formatted_address,
				placeId: place.place_id,
				name: place.name,
			};

			onLocationSelect(location);
			setInputValue(location.address);
		} catch {
			setErrorTips('Error retrieving location data.');
		} finally {
			setIsLoading(false);
		}
	};
	const initAutocomplete = () => {
		if (!inputRef.current?.input) return;

		const options = {
			types: ['(regions)'], // This restricts to cities, states, and zip codes
			componentRestrictions: { country: 'us' }, // Optional: restrict to US
		};

		autocompleteRef.current = new window.google.maps.places.Autocomplete(
			inputRef.current?.input,
			options,
		);

		autocompleteRef.current.addListener('place_changed', handlePlaceSelect);
	};

	const handleInputChange = (e) => {
		setInputValue(e.target.value);
		setIsSelecting(false);
		setErrorTips('');
	};

	// const handleInputBlur = debounce(() => {
	// 	if (inputValue && !isSelecting) {
	// 		setErrorTips('Please select an option after input your location.');
	// 	}
	// }, 1000);

	useEffect(() => {
		if (window.google) {
			initAutocomplete();
		}

		// 保存当前的input引用用于清理
		const currentInputRef = inputRef.current;
		return () => {
			if (autocompleteRef.current) {
				if (currentInputRef?.input) {
					window.google.maps.event.clearInstanceListeners(currentInputRef.input);
				}
				autocompleteRef.current = null;
			}
		};
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	useEffect(() => {
		onError(errorTips);
	}, [errorTips, onError]);

	useImperativeHandle(ref, () => ({
		setErrorTips,
	}));

	return (
		<Input
			ref={inputRef}
			allowClear
			placeholder={placeholder}
			onChange={handleInputChange}
			// onBlur={handleInputBlur}
			value={inputValue}
			size="large"
			prefix={
				<DistanceSvg
					className="common-fill"
					style={{
						color: isSelecting
							? 'var(--base-primary)'
							: errorTips
								? 'var(--base-error)'
								: 'var(--foreground)',
					}}
				/>
			}
			status={errorTips ? 'error' : ''}
			{...props}
		/>
	);
};

export default forwardRef(GeoLocationInput);
