package web

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Windows API constants
const (
	GWL_STYLE     = ^uintptr(16 - 1) // -16 as uintptr
	GWL_EXSTYLE   = ^uintptr(20 - 1) // -20 as uintptr
	WS_BORDER     = 0x00800000
	WS_DLGFRAME   = 0x00400000
	WS_THICKFRAME = 0x00040000
	WS_CAPTION    = 0x00C00000
)

// Windows API structs
type RECT struct {
	Left, Top, Right, Bottom int32
}

type WINDOWPLACEMENT struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
}

type POINT struct {
	X, Y int32
}

// WindowBorderInfo holds border thickness information
type WindowBorderInfo struct {
	Left           int32 `json:"left"`             // Left border (assumed equal to right and bottom)
	Top            int32 `json:"top"`              // Top border (title bar + top frame)
	Right          int32 `json:"right"`            // Right border (equal to left)
	Bottom         int32 `json:"bottom"`           // Bottom border (equal to left/right)
	TitleBarHeight int32 `json:"title_bar_height"` // Title bar height only
	ClientWidth    int32 `json:"client_width"`     // Current client area width
	ClientHeight   int32 `json:"client_height"`    // Current client area height
	WindowWidth    int32 `json:"window_width"`     // Current window width
	WindowHeight   int32 `json:"window_height"`    // Current window height
}

// Windows API functions
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetWindowRect    = user32.NewProc("GetWindowRect")
	procGetClientRect    = user32.NewProc("GetClientRect")
	procGetWindowLongPtr = user32.NewProc("GetWindowLongPtrW")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procFindWindow       = user32.NewProc("FindWindowW")
)

// GetWindowBorderInfo retrieves detailed border information for a window
func GetWindowBorderInfo(hwnd uintptr) (*WindowBorderInfo, error) {
	var windowRect, clientRect RECT

	// Get window rectangle (including borders)
	ret, _, err := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&windowRect)))
	if ret == 0 {
		return nil, fmt.Errorf("GetWindowRect failed: %v", err)
	}

	// Get client rectangle (content area only)
	ret, _, err = procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&clientRect)))
	if ret == 0 {
		return nil, fmt.Errorf("GetClientRect failed: %v", err)
	}

	// Calculate dimensions
	windowWidth := windowRect.Right - windowRect.Left
	windowHeight := windowRect.Bottom - windowRect.Top
	clientWidth := clientRect.Right - clientRect.Left
	clientHeight := clientRect.Bottom - clientRect.Top

	// Get window style to determine border types
	style, _, _ := procGetWindowLongPtr.Call(hwnd, GWL_STYLE)

	// Calculate border sizes
	var borderInfo WindowBorderInfo

	// Use system metrics for standard border sizes
	if style&WS_THICKFRAME != 0 {
		// Resizable window border
		borderInfo.Left = getSystemMetric(7)   // SM_CXBORDER
		borderInfo.Right = getSystemMetric(7)  // SM_CXBORDER
		borderInfo.Top = getSystemMetric(8)    // SM_CYBORDER
		borderInfo.Bottom = getSystemMetric(8) // SM_CYBORDER

		// Add sizing border
		borderInfo.Left += getSystemMetric(32)   // SM_CXSIZEFRAME
		borderInfo.Right += getSystemMetric(32)  // SM_CXSIZEFRAME
		borderInfo.Top += getSystemMetric(33)    // SM_CYSIZEFRAME
		borderInfo.Bottom += getSystemMetric(33) // SM_CYSIZEFRAME
	} else if style&WS_BORDER != 0 || style&WS_DLGFRAME != 0 {
		// Fixed border
		borderInfo.Left = getSystemMetric(5)   // SM_CXFIXEDFRAME
		borderInfo.Right = getSystemMetric(5)  // SM_CXFIXEDFRAME
		borderInfo.Top = getSystemMetric(6)    // SM_CYFIXEDFRAME
		borderInfo.Bottom = getSystemMetric(6) // SM_CYFIXEDFRAME
	}

	// Calculate title bar height
	if style&WS_CAPTION != 0 {
		borderInfo.TitleBarHeight = getSystemMetric(4) // SM_CYCAPTION
		if style&WS_THICKFRAME != 0 {
			borderInfo.TitleBarHeight += getSystemMetric(33) // SM_CYSIZEFRAME
		}
	}

	// Apply custom border calculation logic:
	// Left, Right, Bottom borders are assumed equal
	// Top border = remaining height difference after accounting for 2x bottom border height

	// Calculate uniform side/bottom border
	sideBorder := (windowWidth - clientWidth) / 2
	if sideBorder > 0 {
		borderInfo.Left = sideBorder
		borderInfo.Right = sideBorder
		borderInfo.Bottom = sideBorder
	}

	// Calculate top border (title bar + frame) = total height diff - 2x bottom border
	totalHeightDiff := windowHeight - clientHeight
	if totalHeightDiff > (2 * sideBorder) {
		borderInfo.Top = totalHeightDiff - (2 * sideBorder)
	}

	// Fill in dimensions
	borderInfo.ClientWidth = clientWidth
	borderInfo.ClientHeight = clientHeight
	borderInfo.WindowWidth = windowWidth
	borderInfo.WindowHeight = windowHeight

	return &borderInfo, nil
}

// getSystemMetric retrieves a system metric value
func getSystemMetric(index int) int32 {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(ret)
}

// CalculateWindowSizeForClient calculates the required window size for a desired client area
func CalculateWindowSizeForClient(clientWidth, clientHeight int32, borderInfo *WindowBorderInfo) (windowWidth, windowHeight int32) {
	windowWidth = clientWidth + borderInfo.Left + borderInfo.Right
	windowHeight = clientHeight + borderInfo.Top + borderInfo.Bottom + borderInfo.TitleBarHeight
	return
}

// CalculateClientSizeForWindow calculates the client area size for a given window size
// Uses the custom border logic: window-width - left - right, window-height - top - 2*bottom
func CalculateClientSizeForWindow(windowWidth, windowHeight int32, borderInfo *WindowBorderInfo) (clientWidth, clientHeight int32) {
	clientWidth = windowWidth - borderInfo.Left - borderInfo.Right
	clientHeight = windowHeight - borderInfo.Top - (2 * borderInfo.Bottom)
	return
}

// CalculateWindowPositionForRDP calculates window position and size for RDP winposstr
// Returns: left, top, right, bottom coordinates for the final RDP window
func CalculateWindowPositionForRDP(posX, posY, windowWidth, windowHeight int32) (left, top, right, bottom int32) {
	left = posX
	top = posY
	right = posX + windowWidth
	bottom = posY + windowHeight
	return
}

// GenerateWinPosStr generates the winposstr format for RDP files
// Format: "0,1,<x>,<y>,<right>,<bottom>"
func GenerateWinPosStr(posX, posY, windowWidth, windowHeight int32) string {
	left, top, right, bottom := CalculateWindowPositionForRDP(posX, posY, windowWidth, windowHeight)
	return fmt.Sprintf("0,1,%d,%d,%d,%d", left, top, right, bottom)
}

// findWindowByTitle finds a window by its title
func findWindowByTitle(title string) uintptr {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	return hwnd
}
