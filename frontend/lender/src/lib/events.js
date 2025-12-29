export default class Events {
	constructor() {
		this.events = {};
	}
	on(eventName, callBack) {
		if (!this.events[eventName]) {
			this.events[eventName] = [];
		}
		this.events[eventName].push(callBack);
	}
	emit(eventName, params) {
		if (this.events[eventName]) {
			this.events[eventName].map((callBack) => {
				callBack(params);
			});
		}
	}
	del(eventName, fn) {
		if (this.events[eventName]) {
			this.events[eventName] = this.events[eventName].filter((cb) => fn !== cb);
		}
	}
}
export const generalEvent = new Events(); //通用事件
