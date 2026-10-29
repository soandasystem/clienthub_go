package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"clienthub_go/config"
)

// Estructuras para la API de Resend
type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

type resendResponse struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
}

// sendViaResend envía un email usando la API REST HTTPS de Resend (puerto 443)
func sendViaResend(apiKey, from, toEmail, subject, htmlBody string) error {
	if from == "" {
		from = "TourManager <onboarding@resend.dev>"
	}

	payload := resendRequest{
		From:    from,
		To:      []string{toEmail},
		Subject: subject,
		Html:    htmlBody,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error codificando payload para Resend: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición a Resend: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error conectando a API Resend: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var resErr resendResponse
		if err := json.Unmarshal(bodyBytes, &resErr); err == nil && resErr.Message != "" {
			return fmt.Errorf("error Resend (%d): %s", resp.StatusCode, resErr.Message)
		}
		return fmt.Errorf("error Resend (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("[Resend] Email enviado exitosamente a %s\n", toEmail)
	return nil
}

// getResendKey obtiene la key de Resend desde la config o variables de entorno
func getResendKey(cfg config.Config) string {
	if cfg.ResendAPIKey != "" {
		return cfg.ResendAPIKey
	}
	if key := os.Getenv("PI_KEY_RESEND"); key != "" {
		return key
	}
	if key := os.Getenv("API_KEY_RESEND"); key != "" {
		return key
	}
	if key := os.Getenv("RESEND_API_KEY"); key != "" {
		return key
	}
	return ""
}

// SendPaymentNotification envía un correo de notificación de pago al email indicado.
// Se espera que se llame desde una goroutine para no bloquear el flujo principal.
func SendPaymentNotification(cfg config.Config, toEmail, alumno, estado, monto, commerceOrder string) error {
	if toEmail == "" {
		return fmt.Errorf("email destinatario vacío, no se envía notificación")
	}

	// Determinar asunto y color según estado
	var subject, colorEstado, mensajeEstado string
	switch strings.ToLower(estado) {
	case "pendiente":
		subject = "Pago pendiente de confirmación"
		colorEstado = "#f59e0b"
		mensajeEstado = "Su pago se encuentra pendiente de confirmación. Le notificaremos cuando se complete."
	case "pagado":
		subject = "Pago confirmado exitosamente"
		colorEstado = "#10b981"
		mensajeEstado = "Su pago ha sido procesado y confirmado exitosamente. ¡Gracias por su pago!"
	case "rechazado":
		subject = "Pago rechazado"
		colorEstado = "#ef4444"
		mensajeEstado = "Lamentablemente su pago ha sido rechazado. Por favor intente nuevamente o contacte a soporte."
	case "anulado":
		subject = "Pago anulado"
		colorEstado = "#6b7280"
		mensajeEstado = "Su pago ha sido anulado. Si no realizó esta acción, por favor contacte a soporte."
	default:
		subject = "Notificación de estado de pago"
		colorEstado = "#6b7280"
		mensajeEstado = "El estado de su pago es: " + estado
	}

	// Construir el cuerpo del email en HTML
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 20px;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <div style="background-color: %s; padding: 20px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 24px;">%s</h1>
        </div>
        <div style="padding: 30px;">
            <p style="font-size: 16px; color: #333;">Estimado/a <strong>%s</strong>,</p>
            <p style="font-size: 14px; color: #555;">%s</p>
            <table style="width: 100%%; border-collapse: collapse; margin: 20px 0;">
                <tr style="border-bottom: 1px solid #eee;">
                    <td style="padding: 10px; font-weight: bold; color: #333;">Estado:</td>
                    <td style="padding: 10px; color: %s; font-weight: bold;">%s</td>
                </tr>
                <tr style="border-bottom: 1px solid #eee;">
                    <td style="padding: 10px; font-weight: bold; color: #333;">Monto:</td>
                    <td style="padding: 10px; color: #333;">$%s</td>
                </tr>
                <tr>
                    <td style="padding: 10px; font-weight: bold; color: #333;">Orden:</td>
                    <td style="padding: 10px; color: #333;">%s</td>
                </tr>
            </table>
            <p style="font-size: 12px; color: #999; margin-top: 30px; text-align: center;">
                Este es un correo automático, por favor no responda a este mensaje.
            </p>
        </div>
    </div>
</body>
</html>`, colorEstado, subject, alumno, mensajeEstado, colorEstado, estado, monto, commerceOrder)

	// 1. Si Resend está configurado, usar API HTTP (funciona 100% en Render)
	if resendKey := getResendKey(cfg); resendKey != "" {
		return sendViaResend(resendKey, cfg.ResendFrom, toEmail, subject, body)
	}

	// 2. Fallback a SMTP tradicional
	if cfg.SMTPHost == "" || cfg.SMTPFrom == "" {
		return fmt.Errorf("configuración de email incompleta: configure PI_KEY_RESEND o (SMTP_HOST y SMTP_FROM)")
	}

	// Construir headers del mensaje
	headers := make(map[string]string)
	headers["From"] = cfg.SMTPFrom
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Configurar autenticación SMTP
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	var auth smtp.Auth
	if cfg.SMTPUser != "" && cfg.SMTPPassword != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)
	}

	// Enviar el email
	err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{toEmail}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("error enviando email a %s: %w", toEmail, err)
	}

	fmt.Printf("Email de notificación enviado a %s - Estado: %s\n", toEmail, estado)
	return nil
}

// SendVerificationCodeEmail envía un correo con el código de 6 dígitos para recuperación de contraseña.
func SendVerificationCodeEmail(cfg config.Config, toEmail, code string) error {
	if toEmail == "" {
		return fmt.Errorf("email destinatario vacío, no se envía código")
	}

	if code == "" {
		return fmt.Errorf("código de verificación vacío")
	}

	subject := "Código de recuperación de contraseña"

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f4f6f8; margin: 0; padding: 24px;">
    <div style="max-width: 520px; margin: 0 auto; background-color: #ffffff; border-radius: 10px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); border: 1px solid #e2e8f0;">
        <div style="background: linear-gradient(135deg, #0d6efd 0%%, #0b5ed7 100%%); padding: 24px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 22px; font-weight: 700;">Recuperación de Contraseña</h1>
        </div>
        <div style="padding: 32px 28px;">
            <p style="font-size: 15px; color: #334155; margin-top: 0; line-height: 1.6;">
                Estimado/a usuario/a,
            </p>
            <p style="font-size: 15px; color: #475569; line-height: 1.6;">
                Has solicitado restablecer tu contraseña de acceso. Ingresa el siguiente código de verificación de 6 dígitos para continuar:
            </p>
            <div style="margin: 28px 0; text-align: center;">
                <div style="display: inline-block; background-color: #f8fafc; border: 2px dashed #0d6efd; border-radius: 12px; padding: 16px 36px;">
                    <span style="font-family: 'Courier New', Courier, monospace; font-size: 34px; font-weight: 800; letter-spacing: 8px; color: #0d6efd;">%s</span>
                </div>
            </div>
            <p style="font-size: 13px; color: #64748b; line-height: 1.5; margin-bottom: 0;">
                ⚠️ Este código es confidencial y de uso único. Si no solicitaste este cambio, puedes ignorar este correo de manera segura.
            </p>
            <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 28px 0 20px 0;" />
            <p style="font-size: 12px; color: #94a3b8; margin: 0; text-align: center;">
                Este es un correo automático, por favor no responda a este mensaje.
            </p>
        </div>
    </div>
</body>
</html>`, code)

	// 1. Si Resend está configurado, usar API HTTP (funciona 100% en Render)
	if resendKey := getResendKey(cfg); resendKey != "" {
		return sendViaResend(resendKey, cfg.ResendFrom, toEmail, subject, body)
	}

	// 2. Fallback a SMTP tradicional
	if cfg.SMTPHost == "" || cfg.SMTPFrom == "" {
		return fmt.Errorf("configuración de email incompleta: configure PI_KEY_RESEND o (SMTP_HOST y SMTP_FROM)")
	}

	// Construir headers del mensaje
	headers := make(map[string]string)
	headers["From"] = cfg.SMTPFrom
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Configurar autenticación SMTP
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	var auth smtp.Auth
	if cfg.SMTPUser != "" && cfg.SMTPPassword != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)
	}

	// Enviar el email
	err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{toEmail}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("error enviando código a %s: %w", toEmail, err)
	}

	fmt.Printf("Código de recuperación enviado a %s vía SMTP\n", toEmail)
	return nil
}
