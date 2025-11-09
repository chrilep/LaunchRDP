export namespace main {
	
	export class MonitorWorkArea {
	    index: number;
	    monitorLeft: number;
	    monitorTop: number;
	    monitorRight: number;
	    monitorBottom: number;
	    workLeft: number;
	    workTop: number;
	    workRight: number;
	    workBottom: number;
	    primary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MonitorWorkArea(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.monitorLeft = source["monitorLeft"];
	        this.monitorTop = source["monitorTop"];
	        this.monitorRight = source["monitorRight"];
	        this.monitorBottom = source["monitorBottom"];
	        this.workLeft = source["workLeft"];
	        this.workTop = source["workTop"];
	        this.workRight = source["workRight"];
	        this.workBottom = source["workBottom"];
	        this.primary = source["primary"];
	    }
	}
	export class MousePosition {
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new MousePosition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}
	export class WindowBorderInfo {
	    left: number;
	    right: number;
	    top: number;
	    bottom: number;
	    clientWidth: number;
	    clientHeight: number;
	    windowWidth: number;
	    windowHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowBorderInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.left = source["left"];
	        this.right = source["right"];
	        this.top = source["top"];
	        this.bottom = source["bottom"];
	        this.clientWidth = source["clientWidth"];
	        this.clientHeight = source["clientHeight"];
	        this.windowWidth = source["windowWidth"];
	        this.windowHeight = source["windowHeight"];
	    }
	}
	export class WindowState {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	    deltaX: number;
	    deltaY: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.deltaX = source["deltaX"];
	        this.deltaY = source["deltaY"];
	    }
	}
	export class WorkArea {
	    left: number;
	    top: number;
	    right: number;
	    bottom: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new WorkArea(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.left = source["left"];
	        this.top = source["top"];
	        this.right = source["right"];
	        this.bottom = source["bottom"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}

}

export namespace models {
	
	export class Host {
	    id: string;
	    name: string;
	    address: string;
	    port: number;
	    user_id: string;
	    redirect_clipboard: boolean;
	    redirect_drives: boolean;
	    drives_to_redirect: string;
	    display_mode: string;
	    dynamic_resolution: boolean;
	    screen_mode: number;
	    window_width: number;
	    window_height: number;
	    desktop_width: number;
	    desktop_height: number;
	    position_x: number;
	    position_y: number;
	    win_pos_str: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    modified_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Host(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.user_id = source["user_id"];
	        this.redirect_clipboard = source["redirect_clipboard"];
	        this.redirect_drives = source["redirect_drives"];
	        this.drives_to_redirect = source["drives_to_redirect"];
	        this.display_mode = source["display_mode"];
	        this.dynamic_resolution = source["dynamic_resolution"];
	        this.screen_mode = source["screen_mode"];
	        this.window_width = source["window_width"];
	        this.window_height = source["window_height"];
	        this.desktop_width = source["desktop_width"];
	        this.desktop_height = source["desktop_height"];
	        this.position_x = source["position_x"];
	        this.position_y = source["position_y"];
	        this.win_pos_str = source["win_pos_str"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.modified_at = this.convertValues(source["modified_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class User {
	    id: string;
	    name: string;
	    username: string;
	    login: string;
	    domain: string;
	    encrypted_password: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    modified_at: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.username = source["username"];
	        this.login = source["login"];
	        this.domain = source["domain"];
	        this.encrypted_password = source["encrypted_password"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.modified_at = this.convertValues(source["modified_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

